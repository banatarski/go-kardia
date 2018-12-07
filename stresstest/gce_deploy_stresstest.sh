#!/bin/bash
#
# Create instances on GCE to host Kardia testnet nodes
# Docker images
KARDIA_GO_IMAGE=gcr.io/strategic-ivy-130823/go-kardia:4939014fd1108c53c4660f8913805fa8abaeead3
KARDIA_SCAN_IMAGE=gcr.io/strategic-ivy-130823/kardia-scan:latest


# Number of nodes and prefix name
NODES=3
ETH_NODES=0
NEO_NODES=0
NAME_PREFIX="kardia-stresstest-"

# Generate tx
NUM_TXS=2000
GENERATE_TXS=""
TXS_DELAY=5

# Indexes of main chain validator nodes
MAIN_CHAIN_VAL_INDEXES="1,2,3"

# ports
PORT=3000
RPC_PORT=8545

# Instance specs
ZONE="asia-southeast1-a"
MACHINE_TYPE="n1-standard-1"
IMAGE="cos-stable-70-11021-67-0"

# Index of instance hosting Kardia-scan
KARDIA_SCAN_NODE=3

ENODES=(
     "enode://7a86e2b7628c76fcae76a8b37025cba698a289a44102c5c021594b5c9fce33072ee7ef992f5e018dc44b98fa11fec53824d79015747e8ac474f4ee15b7fbe860@"
	 "enode://660889e39b37ade58f789933954123e56d6498986a0cd9ca63d223e866d5521aaedc9e5298e2f4828a5c90f4c58fb24e19613a462ca0210dd962821794f630f0@"
	 "enode://2e61f57201ec804f9d5298c4665844fd077a2516cd33eccea48f7bdf93de5182da4f57dc7b4d8870e5e291c179c05ff04100718b49184f64a7c0d40cc66343da@"
)



#############################################################################
# cloud_create_instances creates instances with specific specs
#   cloud_create_instances [name_prefix] [num_nodes] [eth_num_nodes]
# Globals:
#   ZONE                  - Zone of the instances to create
#   MACHINE_TYPE          - Specifies the machine type used for the instances
#   IMAGE                 - Boot image for instance
#   ETH_CUSTOM_CPU        - Number of CPUs used on instance hosting Ethereum node
#   ETH_CUSTOM_MEMORY     - Memory size of instance hosting Ethereum node
#   ETH_BOOT_DISK_SIZE    - Disk size of instance hosting Ethereum node
# Arguments:
#   name_prefix           - name prefix of instances
#   num_nodes             - number of instances
#   eth_num_nodes         - number of instances hosting Ethereum node
#   neo_num_nodes         - number of instances hosting Neo node
# Returns:
#   None
#############################################################################
cloud_create_instances() {
  name_prefix=$1
  num_nodes=$2
  eth_num_nodes=$3
  neo_num_nodes=$4
  nodes=()
  ether_dual_nodes=()

  # Create a sequence of node names, used to create instances
  if [ ${num_nodes} -eq 1 ]; then
    nodes=("$name_prefix")
  fi

  next_node=1

  if [ ${num_nodes} -gt 1 ]; then
      # Sequence of Kardia node names, e.g. 9 nodes: ([prefix]007 [prefix]008 [prefix]009 [prefix]010 [prefix]011 [prefix]012 [prefix]013 [prefix]014 [prefix]015 )
      nodes=($(seq -f ${name_prefix}%03g ${next_node} ${num_nodes}))
      echo "=======> Kardia node:" ${nodes[@]}
  fi


  # Specs of instance hosting Kardia node
  args=(
    --machine-type=${MACHINE_TYPE} \
    --zone ${ZONE} \
    --subnet=default \
    --network-tier=PREMIUM \
    --maintenance-policy=MIGRATE \
    --tags=http-server,https-server \
    --image=${IMAGE} \
    --image-project=cos-cloud \
    --boot-disk-size=30GB \
    --boot-disk-type=pd-standard \
    )

  (
    # DEBUG: Print all commands
    set -x
    gcloud compute instances create "${nodes[@]}" "${args[@]}"
  )

}


#############################################################################
# cloud_find_instances finds instances matching the specified pattern
#   cloud_find_instances [filter]
# Globals:
#   none
# Arguments:
#   filter - instance name or prefix
# Returns:
#   string list - list info of instances name:public_ip:private_ip:status
#############################################################################
cloud_find_instances() {
  filter="$1"
  instances=()

  while read -r name public_ip private_ip status; do
    printf "%-30s | public_ip=%-16s private_ip=%s staus=%s\n" "$name" "$public_ip" "$private_ip" "$status"

    instances+=("$name:$public_ip:$private_ip")
  done < <(gcloud compute instances list \
             --filter "$filter" \
             --format 'value(name,networkInterfaces[0].accessConfigs[0].natIP,networkInterfaces[0].networkIP,status)')
}

#############################################################################
# cloud_find_instances finds instances by names matching the specified prefix
#   cloud_find_instances [name_prefix]
# Globals:
#   NAME_PREFIX
# Arguments:
#   name_prefix - instance name prefix
# Returns:
#   string list - list info of instances name:public_ip:private_ip:status
#############################################################################
cloud_find_instances_by_prefix() {
  name_prefix="$1"
  cloud_find_instances "name~^$name_prefix"
}


#############################################################################
# cloud_find_instances finds instances by name
#   cloud_find_instances [name]
# Globals:
#   none
# Arguments:
#   name - instance name
# Returns:
#   string list - list info of instances name:public_ip:private_ip:status
#############################################################################
cloud_find_instance_by_name() {
  name="$1"
  cloud_find_instances "name=$name"
}


#############################################################################
# Implementation flow
#
# 1. Create instances
# 2. Get name and private IP of instances
# 3. Use name and private IP to SSH instances
# 4. Docker pull and run image to start a new node of testnet
#
#############################################################################
# Create instances
echo "Creating instances..."
cloud_create_instances "$NAME_PREFIX" "$NODES" "$ETH_NODES" "$NEO_NODES"

# Get instance info
cloud_find_instances_by_prefix "$NAME_PREFIX"
echo "instances: ${#instances[@]}"

# Store address of instance in format: [enode][private_ip]:[PORT]
addresses=()
public_ips=()
for (( i=0; i<$NODES; i++ ));
do
    # instance[i] is a string "$name:$public_ip:$private_ip"
    # Split instance[i] by ":" to acquire name, public_ip, and private_ip
    # name and private_ip used to SSH to instance
    IFS=: read -r name public_ip private_ip <<< "${instances[$i]}"
    # 10 seed nodes
    if [ $i -lt 10 ] ; then
        addresses+=("${ENODES[$i]}${private_ip}:${PORT}")
    fi
    public_ips+=("http://${public_ip}:${RPC_PORT}")
done

# Concatenate addresses
peers=$(IFS=, ; echo "${addresses[*]}")
rpc_hosts=$(IFS=, ; echo "${public_ips[*]}")

node_index=1
docker_cmd="docker pull ${KARDIA_GO_IMAGE}"
for info in "${instances[@]}";
do
  IFS=: read -r name public_ip private_ip <<< "$info"

  run_cmd=
  if [ "$node_index" -eq 1 ]; then
    GENERATE_TX="--txs --numTxs ${NUM_TXS} --txsDelay ${TXS_DELAY}"
  else
    GENERATE_TX=""
  fi
  run_cmd="docker run -d -p ${PORT}:${PORT} -p ${RPC_PORT}:${RPC_PORT} --name node${node_index} ${KARDIA_GO_IMAGE} --dev --mainChainValIndexes ${MAIN_CHAIN_VAL_INDEXES} --addr :${PORT} --name node${node_index} --rpc --rpcport ${RPC_PORT} --clearDataDir ${GENERATE_TX} --peer ${peers}"
  if [ "$node_index" -eq "$KARDIA_SCAN_NODE" ]; then
    # instance 3 hosts kardia-scan frontend
    # use http://${public_ip}:8080 to see the kardia-scan frontend
    run_cmd="$run_cmd;docker pull ${KARDIA_SCAN_IMAGE}; docker run -d -e "RPC_HOSTS=${rpc_hosts}" -e "publicIP=http://${public_ip}:8080" -p 8080:80 ${KARDIA_SCAN_IMAGE}"
  fi

  # SSH to instance
  (
    echo "Instance $node_index cmds: "

    # DEBUG: Print all commands
    set -x
    gcloud compute ssh "${name}" --zone="${ZONE}" --command="${docker_cmd};${run_cmd}"
  ) &
  ((node_index++))
done
