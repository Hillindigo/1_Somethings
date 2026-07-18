/*
	CommentType          Type = 1
	ApplyType            Type = 2
	DiggArticleType      Type = 3
	UnDiggArticleType    Type = 4
	DiggCommentType      Type = 5
	UnDiggCommentType    Type = 6
	CollectArticleType   Type = 7
	UnCollectArticleType Type = 8
	SystemType           Type = 9
 */
import {type baseResponse, type listResponse, type paramsType, useAxios} from "@/api/index";

export interface siteMsgListType {
    "id": number
    "createdAt": string
    "updatedAt": string
    "type": number
    "revUserID": number
    "ActionUserID": number
    "actionUserNickname": string
    "actionUserAvatar": string
    "title": string
    "content": string
    "articleID": number
    "articleTitle": string
    "commentID": number
    "linkTitle": string
    "linkHref": string
    "isRead": boolean
    "relation"?: 0 | 1 | 2 | 3 | 4
}

export interface siteMsgListParams extends paramsType {
    t: 1 | 2 | 3 // 1 评论和回复 2 赞和收藏 3 系统
}

export function siteMsgListApi(params: siteMsgListParams): Promise<baseResponse<listResponse<siteMsgListType>>> {
    return useAxios.get("/api/site_msg", {params})
}

export interface siteMsgRequest {
    "id": number
    "t": 1 | 2 | 3
}

export function siteMsgReadApi(data: siteMsgRequest): Promise<baseResponse<string>> {
    return useAxios.post("/api/site_msg", data)
}

export function siteMsgRemoveApi(data: siteMsgRequest): Promise<baseResponse<string>> {
    return useAxios.delete("/api/site_msg", {data})
}

export interface siteMsgSettingsType {
    "openCommentMessage": boolean
    "openDiggMessage": boolean
    "openPrivateChat": boolean
}

export function siteMsgSettingsApi(): Promise<baseResponse<siteMsgSettingsType>> {
    return useAxios.get("/api/site_msg/conf")
}

export function siteMsgSettingsUpdateApi(data: siteMsgSettingsType): Promise<baseResponse<string>> {
    return useAxios.put("/api/site_msg/conf", data)
}

export interface unReadMsgType {
    "commentMsgCount": number
    "diggMsgCount": number
    "privateMsgCount": number
    "systemMsgCount": number
}

export function unReadMsgApi(): Promise<baseResponse<unReadMsgType>> {
    return useAxios.get("/api/site_msg/user")
}