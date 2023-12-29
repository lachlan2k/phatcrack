export const credsMap = {
    default: {
        username: 'admin',
        password: 'changeme'
    },

    admin: {
        username: 'admin',
        password: 'd#kdn19SeFiS0@3k'
    },

    alice: {
        username: 'alice',
        password: '84a30c9b39da5951'
    },

    bob: {
        username: 'bobby',
        password: '63e698d946dfbf01'
    },

    carol: {
        username: 'carol',
        password: '9b4b5282dc0cc8fb'
    }
}

const authFileBase = 'playwright/.auth'

export function getAuthFilePath(username: string): string {
    return `${authFileBase}/${username}.json`
}