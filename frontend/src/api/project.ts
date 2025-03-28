import type {
  AttackDTO,
  AttackIDTreeMultipleDTO,
  AttackMultipleDTO,
  AttackStartResponseDTO,
  AttackWithJobsMultipleDTO,
  HashlistAppendRequestDTO,
  HashlistAppendResponseDTO,
  HashlistDTO,
  HashlistResponseMultipleDTO,
  ProjectAddShareRequestDTO,
  ProjectSharesDTO,
  RunningJobCountPerUsersDTO,
  RunningJobsForUserResponseDTO,
  AttackCreateRequestDTO,
  HashlistCreateRequestDTO,
  HashlistCreateResponseDTO,
  ProjectCreateRequestDTO,
  ProjectDTO,
  ProjectResponseMultipleDTO
} from './types'

import { client } from '.'

export function createProject(name: string, description: string): Promise<ProjectDTO> {
  return client
    .post('/api/v1/project/create', {
      name,
      description
    } as ProjectCreateRequestDTO)
    .then(res => res.data)
}

export function deleteProject(projId: string): Promise<string> {
  return client.delete(`/api/v1/project/${projId}`).then(res => res.data)
}

export function getAllProjects(): Promise<ProjectResponseMultipleDTO> {
  return client.get('/api/v1/project/all').then(res => res.data)
}

export function getProject(projId: string): Promise<ProjectDTO> {
  return client.get(`/api/v1/project/${projId}`).then(res => res.data)
}

export function getProjectShares(projId: string): Promise<ProjectSharesDTO> {
  return client.get(`/api/v1/project/${projId}/shares`).then(res => res.data)
}

export function addProjectShare(projId: string, body: ProjectAddShareRequestDTO): Promise<ProjectSharesDTO> {
  return client.post(`/api/v1/project/${projId}/shares`, body).then(res => res.data)
}

export function deleteProjectShare(projId: string, userId: string): Promise<ProjectSharesDTO> {
  return client.delete(`/api/v1/project/${projId}/shares/${userId}`).then(res => res.data)
}

export function createHashlist(body: HashlistCreateRequestDTO): Promise<HashlistCreateResponseDTO> {
  return client.post(`/api/v1/hashlist/create`, body).then(res => res.data)
}

export function appendToHashlist(hashlistId: string, hashes: string[]): Promise<HashlistAppendResponseDTO> {
  return client
    .post(`/api/v1/hashlist/${hashlistId}/append`, {
      input_hashes: hashes
    } as HashlistAppendRequestDTO)
    .then(res => res.data)
}

export function createAttack(body: AttackCreateRequestDTO): Promise<AttackDTO> {
  return client.post(`/api/v1/attack/create`, body).then(res => res.data)
}

export function deleteAttack(attackId: string): Promise<string> {
  return client.delete(`/api/v1/attack/${attackId}`).then(res => res.data)
}

export function stopAttack(attackId: string): Promise<string> {
  return client.delete(`/api/v1/attack/${attackId}/stop`).then(res => res.data)
}

export function startAttack(attackId: string): Promise<AttackStartResponseDTO> {
  return client.put(`/api/v1/attack/${attackId}/start`).then(res => res.data)
}

export function restartAttackFailedJobs(attackId: string): Promise<string> {
  return client.put(`/api/v1/attack/${attackId}/restart-failed-jobs`).then(res => res.data)
}

export function getHashlistsForProject(projId: string): Promise<HashlistResponseMultipleDTO> {
  return client.get(`/api/v1/project/${projId}/hashlists`).then(res => res.data)
}

export function getHashlist(hashlistId: string): Promise<HashlistDTO> {
  return client.get(`/api/v1/hashlist/${hashlistId}`).then(res => res.data)
}

export function deleteHashlist(hashlistId: string): Promise<string> {
  return client.delete(`/api/v1/hashlist/${hashlistId}`).then(res => res.data)
}

export function getAttacksForHashlist(hashlistId: string): Promise<AttackMultipleDTO> {
  return client.get(`/api/v1/hashlist/${hashlistId}/attacks`).then(res => res.data)
}

export function getAttacksWithJobsForHashlist(hashlistId: string, includeRuntimeData: boolean = true): Promise<AttackWithJobsMultipleDTO> {
  return client
    .get(`/api/v1/hashlist/${hashlistId}/attacks-with-jobs` + (includeRuntimeData ? '?includeRuntimeData' : ''))
    .then(res => res.data)
}

export function getAttacksInitialising(): Promise<AttackIDTreeMultipleDTO> {
  return client.get('/api/v1/attack/all-initialising').then(res => res.data)
}

export function getRunningJobs(): Promise<RunningJobsForUserResponseDTO> {
  return client.get('/api/v1/job/all-running').then(res => res.data)
}

export function getJobCountPerUser(): Promise<RunningJobCountPerUsersDTO> {
  return client.get('/api/v1/job/running-count-per-user').then(res => res.data)
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
