import { client } from '.'
import type { ListsGetAllRuleFilesDTO, ListsGetAllWordlistsDTO } from './types'

export function getAllWordlists(): Promise<ListsGetAllWordlistsDTO> {
  return client.get('/api/v1/list/wordlist/all').then((res) => res.data)
}

export function getAllRulefiles(): Promise<ListsGetAllRuleFilesDTO> {
  return client.get('/api/v1/list/rulefile/all').then((res) => res.data)
}
