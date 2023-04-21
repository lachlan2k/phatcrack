import { client } from '.'
import type { AttackDTO, AttackStartResponseDTO } from './types'
import type { AttackCreateRequestDTO } from './types'
import type {
  HashlistCreateRequestDTO,
  HashlistCreateResponseDTO,
  ProjectCreateDTO,
  ProjectDTO,
  ProjectResponseMultipleDTO
} from './types'

export function createProject(name: string, description: string): Promise<ProjectDTO> {
  return client
    .post('/api/v1/project/create', {
      name,
      description
    } as ProjectCreateDTO)
    .then((res) => res.data)
}

export function getAllProjects(): Promise<ProjectResponseMultipleDTO> {
  return client.get('/api/v1/project').then((res) => res.data)
}

export function createHashlist(body: HashlistCreateRequestDTO): Promise<HashlistCreateResponseDTO> {
  return client
    .post(`/api/v1/project/${body.project_id}/hashlist/create`, body)
    .then((res) => res.data)
}

export function createAttack(
  projectId: string,
  hashlistId: string,
  body: AttackCreateRequestDTO
): Promise<AttackDTO> {
  return client
    .post(`/api/v1/hashlist/${hashlistId}/attack/create`)
    .then((res) => res.data)
}

export function startAttack(attackId: string): Promise<AttackStartResponseDTO> {
  return client.put(`/api/v1/attack/${attackId}/start`).then((res) => res.data)
}
