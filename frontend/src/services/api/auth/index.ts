export const utils = {
  placePermission: (place: string, ability: string) => {
    return place + '.' + ability
  },
}

export * from './enums'
export * from './types'
export * from './routes'
