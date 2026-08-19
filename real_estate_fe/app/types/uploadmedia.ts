export interface FileItem {
    id: string
    file: File
    fileType: string;
    previewUrl: string
    thumbnailUrl?: string
    status: 'pending' | 'gettingPresign' | 'uploading' | 'confirming' | 'done' | 'error'
    progress: number
    key?: string
    uploadUrl?: string
    publicUrl?: string
    imageId?: number
    expiresAt?: number
    errorMessage?: string
}

export interface PresignResponse {
    success: boolean
    data?: {
        upload_url: string
        key: string
        expires_at: string
    }
    message?: string
}

export interface ConfirmResponse {
    success: boolean
    data?: {
        image_id: number
        public_url: string
        key: string
        thumbnail_url?: string
    }
    message?: string
}