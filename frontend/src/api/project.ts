import { client } from '.'
import type { ProjectCreateDTO, ProjectDTO, ProjectResponseMultipleDTO } from './types'

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
