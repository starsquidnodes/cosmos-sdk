#!/bin/bash

versions=("v0.47" "v0.50" "latest")

for version in "${versions[@]}"
do
  docker run --rm --platform linux/amd64 ghcr.io/cosmos/simapp:$version bash -c 'echo "movie drastic pumpkin response unhappy morning left deputy genre world margin march wave oven prize sport moment divide spring rare sting coin hockey picture" |
  simd keys add mykey --keyring-backend test --recover > /dev/null 2>&1;
  simd tx staking edit-validator --keyring-backend=test --account-number=1 --sequence=2 --commission-rate=0.25 --offline --sign-mode=amino-json --from=mykey --generate-only > unsigned_tx.json 2>&1;
  simd tx sign --keyring-backend=test unsigned_tx.json --chain-id=testchain --from=mykey --offline --sign-mode=amino-json --account-number=1 --sequence=2' > "signed_$version.tx"
done