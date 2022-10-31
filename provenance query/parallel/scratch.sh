#!/bin/zsh
zmodload zsh/datetime
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

# Environment variables for Org1

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

function getTiming(){
    start=$1
    end=$2

    start_s=$(echo $start | cut -d '.' -f 1)
    start_ns=$(echo $start | cut -d '.' -f 2)
    end_s=$(echo $end | cut -d '.' -f 1)
    end_ns=$(echo $end | cut -d '.' -f 2)


    time=$(( ( 10#$end_s - 10#$start_s ) * 1000 + ( 10#$end_ns / 1000000 - 10#$start_ns / 1000000 ) ))

    echo "$time ms"
}

startTime_s=${epochtime[1]}.${epochtime[2]}
UPC=$1
shift 1
code="code"
for i in "$@"; do
  peer chaincode query -C $i -n $i$code -c '{"function":"IntraChainTransfer","Args":["'$UPC'"]}'>>logfile.txt &
done
wait
sort logfile.txt
endTime_s=${epochtime[1]}.${epochtime[2]}
sumTime=$(getTiming $startTime_s $endTime_s)
echo "$startTime_s ---> $endTime_s" "Total:$sumTime"
rm logfile.txt
echo "end"
