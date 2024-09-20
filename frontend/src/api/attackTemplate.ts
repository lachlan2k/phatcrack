import type {
  AttackTemplateCreateRequestDTO,
  AttackTemplateCreateSetRequestDTO,
  AttackTemplateDTO,
  AttackTemplateGetAllResponseDTO,
  AttackTemplateUpdateRequestDTO
} from './types'

import { client } from '.'

export const AttackTemplateType = 'attack-template'
export const AttackTemplateSetType = 'attack-template-set'

export function getAllAttackTemplates(): Promise<AttackTemplateGetAllResponseDTO> {
  return client.get('/api/v1/attack-template/all').then(res => res.data)
}

export function deleteAttackTemplate(id: string): Promise<string> {
  return client.delete('/api/v1/attack-template/' + id).then(res => res.data)
}

export function createAttackTemplate(newTemplate: AttackTemplateCreateRequestDTO): Promise<AttackTemplateDTO> {
  return client.post('/api/v1/attack-template/create', newTemplate).then(res => res.data)
}

export function createAttackTemplateSet(newTemplateSet: AttackTemplateCreateSetRequestDTO): Promise<AttackTemplateDTO> {
  return client.post('/api/v1/attack-template/create-set', newTemplateSet).then(res => res.data)
}

export function updateAttackTemplate(id: string, body: AttackTemplateUpdateRequestDTO): Promise<AttackTemplateDTO> {
  return client.put('/api/v1/attack-template/' + id, body).then(res => res.data)
}
