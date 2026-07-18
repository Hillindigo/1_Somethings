import type {optionsType} from "@/api";

export interface optionsColorType extends optionsType {
    color?: string
}

export const articleStatusOptions: optionsColorType[] = [
    {label: "草稿", value: 1, color: "green"},
    {label: "审核中", value: 2, color: "red"},
    {label: "已发布", value: 3, color: "blue"},
    {label: "审核拒绝", value: 4, color: "red"},
]

export const roleOptions: optionsColorType[] = [
    {label: "管理员", value: 1, color: "blue"},
    {label: "用户", value: 2, color: "green"},
    {label: "访客", value: 3, color: "red"},
]

export const logLevelOptions: optionsColorType[] = [
    {label: "info", value: 1, color: "blue"},
    {label: "warn", value: 2, color: "orange"},
    {label: "error", value: 3, color: "red"},
]

export const registerSourceOptions: optionsColorType[] = [
    {label: "邮箱注册", value: 1, color: "blue"},
    {label: "QQ注册", value: 2, color: "orange"},
    {label: "命令行注册", value: 3, color: "red"},
]

export const relationOptions: optionsColorType[] = [
    {label: "陌生人", value: 1, color: "red"},
    {label: "已关注", value: 2, color: "orange"},
    {label: "粉丝", value: 3, color: "green"},
    {label: "好友", value: 4, color: "blue"},
]

export const bannerTypeOptions: optionsColorType[] = [
    {label: "banner", value: 1, color: "red"},
    {label: "推广", value: 2, color: "orange"},
]

export const siteMsgTypeOptions: optionsColorType[] = [
    {label: "评论了你的文章", value: 1, color: ""},
    {label: "回复了你的评论", value: 2, color: ""},
    {label: "点赞了你的文章", value: 3, color: ""},
    {label: "取消了点赞", value: 4, color: ""},
    {label: "点赞了你的评论", value: 5, color: ""},
    {label: "取消了评论点赞", value: 6, color: ""},
    {label: "收藏了你的文章", value: 7, color: ""},
    {label: "取消了收藏", value: 8, color: ""},
    {label: "系统通知", value: 9, color: ""},
]