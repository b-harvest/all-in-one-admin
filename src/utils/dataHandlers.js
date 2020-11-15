import { nodeStatusData, chains } from '@config/nodeInfo';
import { getNodeStatus } from '@config/proto3';


export function getIntervalChainData(reduxDataSet) {

    //variable to gather gRPC responses
    let originNetworkData = JSON.parse(JSON.stringify(nodeStatusData))

    // get and set gRPC responses at networkData variable
    setInterval(() => {
        for (const item in chains) {
            chains[item].forEach(info => {
                getNodeStatus(`tcp://${info.ip}:26657`, `${item}/${info.nodeName}`, originNetworkData.main.nodeStatus)
            });
        }
    }, 3000)

    //set redux store data
    setInterval(() => {
        // check node status
        for (let network in originNetworkData.main.nodeStatus) {

            for (let nodeName in originNetworkData.main.nodeStatus[network].nodes) {
                let node = originNetworkData.main.nodeStatus[network].nodes[nodeName]

                if (node.catching_up) {
                    originNetworkData.main.nodeStatus[network].isAllOk = false
                    originNetworkData.main.overview.networks[network].error = true
                } else {
                    originNetworkData.main.nodeStatus[network].isAllOk = true
                    originNetworkData.main.overview.networks[network].error = false
                }
            }
        }

        const newNetworkDataStatus = JSON.parse(JSON.stringify(originNetworkData))
        reduxDataSet(newNetworkDataStatus)
    }, 5000)
}