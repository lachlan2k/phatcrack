import { client } from '.'
import type {
  AttackDTO,
  AttackMultipleDTO,
  AttackStartResponseDTO,
  AttackWithJobsMultipleDTO,
  HashlistDTO,
  HashlistResponseMultipleDTO
} from './types'
import type { AttackCreateRequestDTO } from './types'
import type {
  HashlistCreateRequestDTO,
  HashlistCreateResponseDTO,
  ProjectCreateRequestDTO,
  ProjectDTO,
  ProjectResponseMultipleDTO
} from './types'

export function createProject(name: string, description: string): Promise<ProjectDTO> {
  return client
    .post('/api/v1/project/create', {
      name,
      description
    } as ProjectCreateRequestDTO)
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

export function stopAttack(attackId: string): Promise<string> {
  return client.delete(`/api/v1/attack/${attackId}/stop`).then((res) => res.data)
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

export function getAttacksForHashlist(hashlistId: string): Promise<AttackMultipleDTO> {
  return client.get(`/api/v1/hashlist/${hashlistId}/attacks`).then((res) => res.data)
}

export function getAttacksWithJobsForHashlist(hashlistId: string, includeRuntimeData: boolean = true): Promise<AttackWithJobsMultipleDTO> {
  return client
    .get(`/api/v1/hashlist/${hashlistId}/attacks-with-jobs` + (includeRuntimeData ? '?includeRuntimeData' : ''))
    .then((res) => res.data)
}

export const JobStatusCreated = 'JobStatus-Created'
export const JobStatusAwaitingStart = 'JobStatus-AwaitingStart'
export const JobStatusStarted = 'JobStatus-Started'
export const JobStatusExited = 'JobStatus-Exited'

// Clean exit
export const JobStopReasonFinished = 'JobStopReason-Finished'
// User stopped it
export const JobStopReasonUserStopped = 'JobStopReason-UserStopped'
// Never started in the first place
export const JobStopReasonFailedToStart = 'JobStopReason-FailedToStart'
// General failure
export const JobStopReasonFailed = 'JobStopReason-Failed'
// Agent timed out and we lost contact
export const JobStopReasonTimeout = 'JobStopReason-Timeout'
