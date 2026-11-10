type Filter = object | undefined

export const queryKeys = {
  auth: {
    profile: () => ['auth', 'profile'] as const,
  },

  apps: {
    my: () => ['apps', 'my'] as const,
    adminList: () => ['apps', 'adminList'] as const,
    userAccess: (userId: number) => ['apps', 'userAccess', userId] as const,
  },

  users: {
    all: () => ['users'] as const,
    list: (filter?: Filter) => ['users', 'list', filter] as const,
  },

  roles: {
    options: () => ['roles', 'options'] as const,
  },
}
