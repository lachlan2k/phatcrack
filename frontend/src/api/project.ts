import { client } from '.'
import type { ProjectCreateDTO, ProjectSimpleDetailsDTO, ProjectResponseMultipleDTO } from './types'

export function createProject(name: string, description: string): Promise<ProjectSimpleDetailsDTO> {
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
