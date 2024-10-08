import type { AxiosProgressEvent } from 'axios'

import type { GetAllListfilesDTO, ListfileDTO } from './types'

import { client } from '.'

export const LISTFILE_TYPE_WORDLIST = 'Wordlist'
export const LISTFILE_TYPE_RULEFILE = 'Rulefile'

export type ListfileTypeT = 'Wordlist' | 'Rulefile'

export function getAllListfiles(): Promise<GetAllListfilesDTO> {
  return client.get('/api/v1/listfiles/all').then(res => res.data)
}

export function getListfilesForProject(projectId: string): Promise<GetAllListfilesDTO> {
  return client.get(`/api/v1/project/${projectId}/listfiles`).then(res => res.data)
}

export function deleteListfile(id: string): Promise<string> {
  return client.delete('/api/v1/listfiles/' + id).then(res => res.data)
}

export function uploadListfile(body: FormData, onUploadProgress: (progress: AxiosProgressEvent) => void): Promise<ListfileDTO> {
  return client
    .post('/api/v1/listfiles/upload', body, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress
    })
    .then(res => res.data)
}
