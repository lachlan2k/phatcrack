import * as api from '../../../frontend/src/api'
import {
  beforeAllSetupClientWithCookieJar,
  beforeAllSetupClientWithLogin,
  enrichAxiosError,
  setupClientWithCookieJar,
  initialAdminPassword,
  credsMap,
  re400,
  re401,
  re40x,
  reOk
} from './_helpers'

describe('Projects', () => {
  let alicesProjectName = 'alices project'
  let alicesProjectId = ''

  describe('Alice creating projects', () => {
    const u = credsMap.alice
    beforeAllSetupClientWithLogin(u)

    it('allows alice to create a project', async () => {
      const res = await api.createProject(alicesProjectName, '')
      expect(res.name).toBe(alicesProjectName)
      alicesProjectId = res.id

      const aliceUserData = await api.refreshAuth()
      expect(res.owner_user_id).toBe(aliceUserData.user.id)
    })

    it('doesnt allow a project with a bad name to be created', async () => {
      const badNames = ['', 'a', 'a'.repeat(1000), '!@#$%^&*()-=_+"\'']

      for (const name of badNames) {
        await expect(api.createProject(name, '')).rejects.toThrow(re400)
      }
    })
  })

  describe('Access Control', () => {
    describe('Alice Perspective', () => {
      beforeAllSetupClientWithLogin(credsMap.alice)

      it('allows alice to see her project in list of all projects', async () => {
        const res = await api.getAllProjects()
        const ids = res.projects.map(x => x.id)

        expect(ids).toContain(alicesProjectId)
      })

      it('allows alice to get her project', async () => {
        const res = await api.getProject(alicesProjectId)
        expect(res.id).toBe(alicesProjectId)
        expect(res.name).toBe(alicesProjectName)
      })
    })

    describe('Bob Perspective', () => {
      beforeAllSetupClientWithLogin(credsMap.bob)

      it('doesnt include alices project in bobs list of all projects', async () => {
        const res = await api.getAllProjects()
        const ids = res.projects.map(x => x.id)

        expect(ids).not.toContain(alicesProjectId)
      })

      it('doesnt allow bob to see alices project', async () => {
        await expect(api.getProject(alicesProjectId)).rejects.toThrow(re40x)
      })

      it('doesnt allow bob to delete alices project', async () => {
        await expect(api.deleteProject(alicesProjectId)).rejects.toThrow(re40x)
      })
    })

    describe('Admin Perspective', () => {
      beforeAllSetupClientWithLogin(credsMap.admin)

      it("allows admin to see alice's project in list of all projects", async () => {
        const res = await api.getAllProjects()
        const ids = res.projects.map(x => x.id)

        expect(ids).toContain(alicesProjectId)
      })

      it("allows admin to get alice's project", async () => {
        const res = await api.getProject(alicesProjectId)
        expect(res.id).toBe(alicesProjectId)
        expect(res.name).toBe(alicesProjectName)
      })
    })
  })
})
