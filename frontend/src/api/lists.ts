import { client } from '.'
import type { GetAllWordlistsDTO, GetAllRuleFilesDTO } from './types'

export function getAllWordlists(): Promise<GetAllWordlistsDTO> {
  return client.get('/api/v1/list/wordlist/all').then((res) => res.data)
}

export function getAllRulefiles(): Promise<GetAllRuleFilesDTO> {
  return client.get('/api/v1/list/rulefile/all').then((res) => res.data)
}
