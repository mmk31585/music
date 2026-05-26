export const Permissions = {
  READ_HISTORY: 'read_history',
  READ: 'read',
  READ_LIST: 'read_list',
  READ_ONE: 'read_one',
  CREATE: 'create',
  UPDATE: 'update',
  DELETE: 'delete',
  PERMANENT_DELETE: 'permanent_delete',
  RESTORE: 'restore',
  PRINT: 'print',
  EXPORT: 'export',
  IMPORT: 'import',
  WATCH: 'watch', // Can user follow any happened interaction
  ARCHIVE: 'archive',
  PUBLISH: 'publish',
  FINALIZE: 'finalize',
  VERIFY: 'verify',
  CHANGE_STATUS: 'change_status',
  CHANGE_PASSWORD: 'change_password',
  DOWNLOAD: 'download',
  UPLOAD: 'upload',
  START: 'start',
  END: 'end',
  BAN: 'ban',
  REPLY: 'reply',
}

export const PermissionPlaces = {}
export enum AuthApiRoutes {
  REGISTER = '/auth/register',
  LOGIN = '/auth/login',
  REFRESH = '/auth/refresh',
  LOGOUT = '/auth/logout',
  ME = '/auth/me',
}
