import type { AttackTemplateCreateRequestDTO, AttackTemplateDTO, AttackTemplateGetAllResponseDTO } from './types'

import { client } from '.'

export function getAllAttackTemplates(): Promise<AttackTemplateGetAllResponseDTO> {
  return client.get('/api/v1/attack-template/all').then(res => res.data)
}

export function deleteAttackTemplate(id: string): Promise<string> {
  return client.delete('/api/v1/attack-template/' + id).then(res => res.data)
}

export function createAttackTemplate(newTemplate: AttackTemplateCreateRequestDTO): Promise<AttackTemplateDTO> {
  return client.post('/api/v1/attack-template/create', newTemplate).then(res => res.data)
}
