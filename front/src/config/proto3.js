import { apiAddress } from './nodeInfo'
import Axios from "axios"

export function getNodeStatus(uri, nodeTag, data) {
    const requestTime = new Date()
    Axios.get(`${process.env.NODE_ENV === "test" ? apiAddress.dev : apiAddress.prod}/GetnodeStatus?nodeuri=${uri}`)
        .then(response => {
            const responseTime = new Date()
            const timeSpent = responseTime - requestTime
            let parsedArray = JSON.parse([response.data.status])
            parsedArray.moniker = `${nodeTag?.split('/')[1]}\n(${timeSpent}ms)`

            let monikerReplacedResponse = parsedArray
            data[nodeTag?.split('/')[0]].nodes[nodeTag?.split('/')[1]] = monikerReplacedResponse
        }).catch(err => {
            const responseTime = new Date()
            const timeSpent = responseTime - requestTime
            data[nodeTag?.split('/')[0]].nodes[nodeTag?.split('/')[1]] = { "latest_block_height": err.message, "catching_up": 'code:' + err.code, "moniker": `${nodeTag?.split('/')[1]}\n(${timeSpent}ms)` }
            console.error(err)
        })
}

export function getValidatorSignInfo(uri, vali_address, nodeTag, data) {
    Axios.get(`${process.env.NODE_ENV === "test" ? apiAddress.dev : apiAddress.prod}/GetvalidatorSignInfo?nodeuri=${uri}&validatoraddress=${vali_address}`)
        .then(response => {
            let res = JSON.parse(response.data.status)
            data[nodeTag?.split('/')[0]].isSign = res.SignInfo
        }).catch(err => {
            console.error(err)
        })
}