#!/bin/bash
#
# Copyright The Linux Foundation All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
set -ex

WORKPATH=/home/zhangyiguang/fabric/fabric-samples

# Bring the test network down
pushd $WORKPATH/test-network
./network.sh down
popd

pushd $WORKPATH/test-network/prometheus-grafana
docker-compose down
popd

# clean out any old identites in the wallets
rm -rf ./wallet/*
rm -rf ./data/ledger.data