import type { AxiosProgressEvent } from 'axios'
import { client } from '.'
import type { GetAllWordlistsDTO, GetAllRuleFilesDTO, ListfileDTO } from './types'

export function getAllWordlists(): Promise<GetAllWordlistsDTO> {
  return client.get('/api/v1/list/wordlist/all').then((res) => res.data)
}

export function getAllRulefiles(): Promise<GetAllRuleFilesDTO> {
  return client.get('/api/v1/list/rulefile/all').then((res) => res.data)
}

export function deleteListfile(id: string): Promise<string> {
  return client.delete('/api/v1/list/listfile/' + id).then((res) => res.data)
}

export function uploadListfile(body: FormData, onUploadProgress: (progress: AxiosProgressEvent) => void): Promise<ListfileDTO> {
  return client
    .post('/api/v1/list/upload', body, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress
    })
    .then((res) => res.data)
}
