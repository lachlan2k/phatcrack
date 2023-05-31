import { client } from '.'
import type {
  AttackDTO,
  AttackStartResponseDTO,
  HashlistDTO,
  HashlistResponseMultipleDTO
} from './types'
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
  return client.get('/api/v1/project/all').then((res) => res.data)
}

export function getProject(projId: string): Promise<ProjectDTO> {
  return client.get(`/api/v1/project/${projId}`).then((res) => res.data)
}

export function createHashlist(body: HashlistCreateRequestDTO): Promise<HashlistCreateResponseDTO> {
  return client.post(`/api/v1/hashlist/create`, body).then((res) => res.data)
}

export function createAttack(body: AttackCreateRequestDTO): Promise<AttackDTO> {
  return client.post(`/api/v1/attack/create`, body).then((res) => res.data)
}

export function startAttack(attackId: string): Promise<AttackStartResponseDTO> {
  return client.put(`/api/v1/attack/${attackId}/start`).then((res) => res.data)
}

export function getHashlistsForProject(projId: string): Promise<HashlistResponseMultipleDTO> {
  return client.get(`/api/v1/project/${projId}/hashlists`).then((res) => res.data)
}

export function getHashlist(hashlistId: string): Promise<HashlistDTO> {
  return client.get(`/api/v1/hashlist/${hashlistId}`).then((res) => res.data)
}
