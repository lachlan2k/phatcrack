import * as api from '../../../frontend/src/api'
import { dummyApiRequests } from './dummyRequests'
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

describe('Hashlists', () => {
  const alicesHashlistName = 'alices hashlist'
  const alicesProjectName = 'project for hashlist testing'
  let alicesProjectId = ''
  let alicesHashlistId = ''
  let alicesHashlistWithUsernamesId = ''

  describe('Alice creating hashlists', () => {
    const u = credsMap.alice
    beforeAllSetupClientWithLogin(u)

    it('allows alice to create a project', async () => {
      const res = await api.createProject(alicesProjectName, '')
      expect(res.name).toBe(alicesProjectName)
      alicesProjectId = res.id
    })

    it('doesnt allow alice to create a hashlist with a bad name', async () => {
      const invalidNames = ['!@#$%^&*()-=_+"\'', '', ' ', '   ', 'a', 'a'.repeat(256), '\n\t\r', '😀🚀🌈']

      for (const invalidName of invalidNames) {
        await expect(
          api.createHashlist({
            name: invalidName,
            project_id: alicesProjectId,
            has_usernames: false,
            hash_type: 0,
            input_hashes: []
          })
        ).rejects.toThrow(re400)
      }
    })

    describe('Hash validation', () => {
      it('allows alice to create a hashlist with valid MD5 hashes', async () => {
        const validMD5Hashes = ['098f6bcd4621d373cade4e832627b4f6', '5d41402abc4b2a76b9719d911017c592', 'd41d8cd98f00b204e9800998ecf8427e']

        const res = await api.createHashlist({
          name: alicesHashlistName,
          project_id: alicesProjectId,
          has_usernames: false,
          hash_type: 0, // Assuming 0 represents MD5
          input_hashes: validMD5Hashes
        })

        expect(res.id).not.toBeNull()
        alicesHashlistId = res.id

        const hashlist = await api.getHashlist(res.id)
        expect(hashlist.hashes.length).toBe(validMD5Hashes.length)
      })

      it('does not allow alice to create a hashlist with invalid MD5 hashes', async () => {
        const invalidMD5Hashes = [
          '098f6bcd4621d373cade4e832627b4f', // Too short
          '5d41402abc4b2a76b9719d911017c592g', // Invalid character
          'notahashatalljusttext'
        ]

        await expect(
          api.createHashlist({
            name: 'Invalid MD5 Hashlist',
            project_id: alicesProjectId,
            has_usernames: false,
            hash_type: 0,
            input_hashes: invalidMD5Hashes
          })
        ).rejects.toThrow(re400)
      })

      it('allows alice to create a hashlist with valid MD5 hashes and usernames', async () => {
        const validMD5HashesWithUsernames = [
          'user1:098f6bcd4621d373cade4e832627b4f6',
          'user2:5d41402abc4b2a76b9719d911017c592',
          'user3:d41d8cd98f00b204e9800998ecf8427e'
        ]

        const res = await api.createHashlist({
          name: 'Valid MD5 Hashlist with Usernames',
          project_id: alicesProjectId,
          has_usernames: true,
          hash_type: 0, // Assuming 0 represents MD5
          input_hashes: validMD5HashesWithUsernames
        })

        expect(res.id).not.toBeNull()
        const hashlistId = res.id
        alicesHashlistWithUsernamesId = hashlistId

        const hashlist = await api.getHashlist(hashlistId)
        expect(hashlist.hashes.length).toBe(validMD5HashesWithUsernames.length)
        expect(hashlist.has_usernames).toBe(true)

        // Verify that usernames are correctly stored
        for (let i = 0; i < validMD5HashesWithUsernames.length; i++) {
          const [username, hash] = validMD5HashesWithUsernames[i].split(':')
          expect(hashlist.hashes[i].username).toBe(username)
          expect(hashlist.hashes[i].input_hash).toBe(hash)
        }
      })

      it('does not allow alice to create a hashlist with invalid MD5 hashes and usernames', async () => {
        const invalidMD5HashesWithUsernames = [
          'user1:098f6bcd4621d373cade4e832627b4f', // Too short
          'user2:5d41402abc4b2a76b9719d911017c592g', // Invalid character
          'user3:notahashatalljusttext'
        ]

        await expect(
          api.createHashlist({
            name: 'Invalid MD5 Hashlist with Usernames',
            project_id: alicesProjectId,
            has_usernames: true,
            hash_type: 0,
            input_hashes: invalidMD5HashesWithUsernames
          })
        ).rejects.toThrow(re400)
      })

      it('does not allow alice to create a hashlist with usernames when has_usernames is false', async () => {
        const hashesWithUsernames = [
          'user1:098f6bcd4621d373cade4e832627b4f6',
          'user2:5d41402abc4b2a76b9719d911017c592',
          'user3:d41d8cd98f00b204e9800998ecf8427e'
        ]

        await expect(
          api.createHashlist({
            name: 'Hashlist with Usernames but has_usernames false',
            project_id: alicesProjectId,
            has_usernames: false,
            hash_type: 0, // MD5
            input_hashes: hashesWithUsernames
          })
        ).rejects.toThrow(re400)
      })

      it('allows alice to create a hashlist with valid bcrypt hashes', async () => {
        const validBcryptHashes = [
          '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
          '$2b$10$2BpL3Ll5CYHePRnYJD7xTOXsru65GFHGxus9tQOdqOXPwWQNRyry.',
          '$2y$10$6bNw2HLQYeqHYyBfLMsv/OiwqTymGIGzFsA4hOTWebfehV7Cj4qHu'
        ]

        const res = await api.createHashlist({
          name: 'Valid Bcrypt Hashlist',
          project_id: alicesProjectId,
          has_usernames: false,
          hash_type: 3200, // bcrypt
          input_hashes: validBcryptHashes
        })

        expect(res.id).toBeDefined()
        const hashlist = await api.getHashlist(res.id)
        expect(hashlist.hashes.length).toBe(validBcryptHashes.length)
        expect(hashlist.has_usernames).toBe(false)
        expect(hashlist.hash_type).toBe(3200)

        for (let i = 0; i < validBcryptHashes.length; i++) {
          expect(hashlist.hashes[i].input_hash).toBe(validBcryptHashes[i])
        }
      })

      it('does not allow alice to create a hashlist with invalid bcrypt hashes', async () => {
        const invalidBcryptHashes = [
          '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhW', // Too short
          '$2b$10$2BpL3Ll5CYHePRnYJD7xTOXsru65GFHGxus9tQOdqOXPwWQNRyry.X', // Too long
          '$2c$10$6bNw2HLQYeqHYyBfLMsv/OiwqTymGIGzFsA4hOTWebfehV7Cj4qHu', // Invalid version
          '$2y$11$6bNw2HLQYeqHYyBfLMsv/OiwqTymGIGzFsA4hOTWebfehV7Cj4qHu', // Invalid cost factor
          'notahashatalljusttext'
        ]

        await expect(
          api.createHashlist({
            name: 'Invalid Bcrypt Hashlist',
            project_id: alicesProjectId,
            has_usernames: false,
            hash_type: 3200, // bcrypt
            input_hashes: invalidBcryptHashes
          })
        ).rejects.toThrow(re400)
      })
    })

    describe('Appending Hashes', () => {
      beforeAllSetupClientWithLogin(credsMap.alice)

      it('allows alice to append valid hashes to her hashlist', async () => {
        const existingHashesToCheckDedup = ['098f6bcd4621d373cade4e832627b4f6']
        const newHashes = ['5f4dcc3b5aa765d61d8327deb882cf99', 'a41d8cd98f00b204e9800998ecf8427f']
        const startingHashlist = await api.getHashlist(alicesHashlistId)
        const startingHashCount = startingHashlist.hashes.length

        const appendResponse = await api.appendToHashlist(alicesHashlistId, [...existingHashesToCheckDedup, ...newHashes])
        expect(appendResponse.num_new_hashes).toBe(newHashes.length)

        const updatedHashlist = await api.getHashlist(alicesHashlistId)
        const appendedHashes = updatedHashlist.hashes.slice(-newHashes.length)

        for (let i = 0; i < newHashes.length; i++) {
          expect(appendedHashes[i].input_hash).toBe(newHashes[i])
        }

        // sanity check the final length too, so it doesn't include the deduped hashes
        expect(updatedHashlist.hashes.length).toBe(startingHashCount + newHashes.length)
      })

      it('allows alice to append valid hashes with usernames to her hashlist with usernames', async () => {
        const newHashesWithUsernames = ['user3:5f4dcc3b5aa765d61d8327deb882cf99', 'user4:098f6bcd4621d373cade4e832627b4f6']

        const appendResponse = await api.appendToHashlist(alicesHashlistWithUsernamesId, newHashesWithUsernames)
        expect(appendResponse.num_new_hashes).toBe(newHashesWithUsernames.length)

        const updatedHashlist = await api.getHashlist(alicesHashlistWithUsernamesId)
        const appendedHashes = updatedHashlist.hashes.slice(-newHashesWithUsernames.length)

        for (let i = 0; i < newHashesWithUsernames.length; i++) {
          const [username, hash] = newHashesWithUsernames[i].split(':')
          expect(appendedHashes[i].username).toBe(username)
          expect(appendedHashes[i].input_hash).toBe(hash)
        }
      })

      it('does not allow alice to append hashes without usernames to her hashlist with usernames', async () => {
        const invalidHashes = ['5f4dcc3b5aa765d61d8327deb882cf99', '098f6bcd4621d373cade4e832627b4f6']

        await expect(api.appendToHashlist(alicesHashlistWithUsernamesId, invalidHashes)).rejects.toThrow(re400)
      })

      it('does not allow alice to append hashes with usernames to her hashlist without usernames', async () => {
        const hashesWithUsernames = ['user1:5f4dcc3b5aa765d61d8327deb882cf99', 'user2:098f6bcd4621d373cade4e832627b4f6']

        await expect(api.appendToHashlist(alicesHashlistId, hashesWithUsernames)).rejects.toThrow(re400)

        // Verify that the hashlist remains unchanged
        const hashlist = await api.getHashlist(alicesHashlistId)
        for (const hash of hashlist.hashes) {
          expect(hash.input_hash).not.toContain(':')
          expect(hash.username).toBe('')
        }
      })

      it('does not allow alice to append invalid hashes to her hashlist', async () => {
        const invalidHashes = [
          'not_a_valid_hash',
          '123' // Too short
        ]

        await expect(api.appendToHashlist(alicesHashlistId, invalidHashes)).rejects.toThrow(re400)
      })
    })
  })

  describe('Access Control', () => {
    describe('Alice Perspective', () => {
      beforeAllSetupClientWithLogin(credsMap.alice)

      it('allows alice to get her hashlist', async () => {
        const res = await api.getHashlist(alicesHashlistId)
        expect(res.id).toBe(alicesHashlistId)
        expect(res.name).toBe(alicesHashlistName)
      })

      it('allows alice to see her hashlist in her project', async () => {
        const res = await api.getHashlistsForProject(alicesProjectId)
        const ids = res.hashlists.map(x => x.id)
        expect(ids).toContain(alicesHashlistId)
      })
    })

    describe('Bob Perspective', () => {
      beforeAllSetupClientWithLogin(credsMap.bob)

      it('doesnt allow bob to get alices hashlist', async () => {
        await expect(api.getHashlist(alicesHashlistId)).rejects.toThrow(re40x)
      })

      it('doesnt allow bob to see alices hashlist in her project', async () => {
        await expect(api.getHashlistsForProject(alicesProjectId)).rejects.toThrow(re40x)
      })

      it('doesnt allow bob to delete alices hashlist', async () => {
        await expect(api.deleteHashlist(alicesHashlistId)).rejects.toThrow(re40x)
      })

      it('doesnt allow bob to append hashes to alices hashlist', async () => {
        const validHashes = ['098f6bcd4621d373cade4e832627b4f6', '5d41402abc4b2a76b9719d911017c592']
        await expect(api.appendToHashlist(alicesHashlistId, validHashes)).rejects.toThrow(re40x)
      })
    })

    describe('Admin Perspective', () => {
      beforeAllSetupClientWithLogin(credsMap.admin)

      it('allows admin to get alices hashlist', async () => {
        const res = await api.getHashlist(alicesHashlistId)
        expect(res.id).toBe(alicesHashlistId)
        expect(res.name).toBe(alicesHashlistName)
      })

      it('allows admin to see alices hashlist in her project', async () => {
        const res = await api.getHashlistsForProject(alicesProjectId)
        const ids = res.hashlists.map(x => x.id)
        expect(ids).toContain(alicesHashlistId)
      })
    })
  })
})
