type Article = {
    id: number
    title: string
    content: string
    likeCnt: number
    liked: boolean
    collectCnt: number
    collected: boolean
    readCnt: number
}

export interface ArticleItem {
    id: bigint;
    title: string;
    status: number;
    abstract: string;
    content?: string;
    authorId?: bigint;
    createdAt?: string;
    updatedAt?: string;
}

export interface ArticleDetail extends ArticleItem {
    readCnt: number;
    likeCnt: number;
    collectCnt: number;
    liked: boolean;
    collected: boolean;
}

export interface ArticleListRequest {
    offset: number;
    limit: number;
}

export interface ArticleListResponse {
    code: number;
    msg: string;
    data: ArticleItem[];
}

export interface ArticleDetailResponse {
    code: number;
    msg: string;
    data: ArticleDetail;
}

export interface ArticleCreateRequest {
    title: string;
    content: string;
    abstract: string;
    status: number;
}

export interface ArticleUpdateRequest extends ArticleCreateRequest {
    id: bigint;
}

export interface ArticleResponse {
    code: number;
    msg: string;
    data?: any;
}

export interface CodeURL {
    codeURL: string;
    rid: number;
}

export interface RewardRequest {
    id: number;
    amt: number;
}

export interface RewardResponse {
    code: number;
    msg: string;
    data: CodeURL;
}

export interface RewardDetailRequest {
    rid: number;
}

export interface RewardDetailResponse {
    code: number;
    msg: string;
    data: string;
}