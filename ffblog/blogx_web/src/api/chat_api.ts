import {type baseResponse, type listResponse, type paramsType, useAxios} from "@/api/index";

export interface textMsgType {
    content: string
}

export interface imageMsgType {
    src: string
}

export interface markdownMsgType {
    content: string
}

export interface readMsgType {
    readChatID: number
}

export interface msgType {
    textMsg?: textMsgType
    imageMsg?: imageMsgType
    markdownMsg?: markdownMsgType
    msgReadMsg?: readMsgType
}

export interface sessionListType {
    "userID": number
    "userNickname": string
    "userAvatar": string
    "msg": msgType
    "msgType": number // 1 2 3  11
    "newMsgDate": string
    "relation": number
}

export interface sessionListParams extends paramsType {

}

export function sessionListApi(params?: paramsType): Promise<baseResponse<listResponse<sessionListType>>> {
    return useAxios.get("/api/chat/session", {params})
}

export interface chatListType {
    "id": number
    "createdAt": string
    "updatedAt": string
    "sendUserID": number
    "revUserID": number
    "msgType": number
    "msg": msgType
    "sendUserNickname": string
    "sendUserAvatar": string
    "revUserNickname": string
    "revUserAvatar": string
    "isMe": boolean
    "isRead": boolean
}

export interface chatListParams extends paramsType {
    revUserID: number
    type: 1 | 2
}

export function chatListApi(params: chatListParams): Promise<baseResponse<listResponse<chatListType>>> {
    return useAxios.get("/api/chat", {params})
}

export interface chatSendType {
    "revUserID": number
    "msgType": number
    "msg": msgType
}

export function chatRemoveApi(idList: number[]):Promise<baseResponse<string>>{
    return useAxios.delete("/api/chat", {data: {idList}})
}
export function sessionRemoveApi(userID:number):Promise<baseResponse<string>>{
    return useAxios.delete("/api/chat/user/" + userID.toString())
}