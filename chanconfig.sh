#!/bin/sh

# 检查是否提供了足够的参数
if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <blockzise> <batch_timeout>"
    exit 1
fi

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_LOCALMSPID="Org1MSP"

new_blocksize=$1
new_batch_timeout=$2
ORDERER_CONTAINER=orderer.example.com:7050
TLS_ROOT_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem

peer channel fetch config config_block.pb -o orderer.example.com:7050 -c mychannel --tls --cafile  /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq .data.data[0].payload.data.config config_block.json > config.json
cp config.json modified_config.json

# 检查输入是否为正整数
# if [[ "$new_blocksize" =~ ^[1-9][0-9]*$ ]]; then
if [ "$new_blocksize" -ge 1 ] 2>/dev/null; then
        old_bs=$(grep max_message_count modified_config.json | awk '{print $2}')
        old_bsconfig="\"max_message_count\": $old_bs"
        new_bsconfig="\"max_message_count\": $new_blocksize,"
        echo "Change blocksize from $old_bs to $new_blocksize"
        sed -i "s/$old_bsconfig/$new_bsconfig/g" modified_config.json
else
        echo "Keep blocksize"
fi
        old_bto=$(grep timeout modified_config.json | awk '{print $2}')
        old_btoconfig="\"timeout\": $old_bto"
        new_btoconfig="\"timeout\": \"$new_batch_timeout\\s\""
        echo "Change timeout from $old_bto to $new_batch_timeout"
        sed -i "s/$old_btoconfig/$new_btoconfig/g" modified_config.json


configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id mychannel --original config.pb --updated modified_config.pb --output config_update.pb
configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"mychannel", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

CORE_PEER_LOCALMSPID=OrdererMSP
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp/
peer channel update -f config_update_in_envelope.pb -c mychannel -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem