// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface QueryParamResult<F = Record<string, any>> {
  page?: number
  per_page?: number
  limit?: number
  sort?: string | string[]
  filter?: F
  include?: string | string[]
  append?: string | string[]
}

class QueryBuilder {
  _path: string
  _params: Record<string, string[]> = {}
  static aliases: Record<string, string> = {}

  constructor(_path: string) {
    this._path = _path
  }

  static defineAliases(aliases: Record<string, string>) {
    QueryBuilder.aliases = aliases
  }

  filter(key: string, value: string | number | boolean) {
    this.createOrUpdateParams(`filter[${key}]`, [String(value)])

    return this
  }

  sort(...sorts: string[]) {
    this.createOrUpdateParams('sort', sorts)

    return this
  }

  include(...includes: string[]) {
    this.createOrUpdateParams('include', includes)

    return this
  }

  append(...appends: string[]) {
    this.createOrUpdateParams('append', appends)

    return this
  }

  param(key: string, ...values: string[]) {
    this.createOrUpdateParams(key, values)

    return this
  }

  fields(data: Record<string, string[]>) {
    Object.entries(data).forEach(([key, values]) => {
      this.createOrUpdateParams(`fields[${key}]`, values)
    })

    return this
  }

  page(page: number) {
    this._params['page'] = [page.toString()]
    return this
  }

  perPage(perPage: number) {
    this._params['per_page'] = [perPage.toString()]
    return this
  }

  limit(limit: number) {
    this._params['limit'] = [limit.toString()]
    return this
  }

  when(condition: boolean, callback: (query: QueryBuilder) => void) {
    if (condition) {
      callback(this)
    }

    return this
  }

  forgetValue(key: string, value: string) {
    if (!this._params[key]) {
      key = this.parseKeyToAlias(key)
    }

    this._params[key] = (this._params[key] || []).filter((val) => val !== value)

    return this
  }

  forget(key: string) {
    if (!this._params[key]) {
      key = this.parseKeyToAlias(key)
    }

    delete this._params[key]

    return this
  }

  forgets(...keys: string[]) {
    keys.forEach((key) => {
      this.forget(key)
    })

    return this
  }

  createOrUpdateParams(key: string, values: string[]) {
    key = this.parseKeyToAlias(key)
    this._params[key] = [...(this._params[key] || []), ...values]

    return this
  }

  parseKeyToAlias(key: string) {
    Object.entries(QueryBuilder.aliases).forEach(([alias, keyAlias]) => {
      if (key.includes(alias)) {
        key = key.replace(alias, keyAlias)
      }
    })

    return key
  }

  tap(callback: (query: QueryBuilder) => void) {
    callback(this)

    return this
  }

  scopes(...scopes: ((query: QueryBuilder) => QueryBuilder)[]) {
    scopes.forEach((scope) => scope(this))

    return this
  }

  getParams() {
    return this._params
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  fromQueryParamResult<F extends Record<string, any>>(result: QueryParamResult<F>) {
    if (!result) return this

    // pagination
    if (result.page !== undefined) this.page(result.page)
    if (result.per_page !== undefined) this.perPage(result.per_page)
    if (result.limit !== undefined) this.limit(result.limit)

    // helper: normalize string | string[] -> string[]
    const asArray = (v?: string | string[]) =>
      v === undefined ? undefined : Array.isArray(v) ? v : [v]

    const sortArr = asArray(result.sort)
    if (sortArr) this.sort(...sortArr)

    const includeArr = asArray(result.include)
    if (includeArr) this.include(...includeArr)

    const appendArr = asArray(result.append)
    if (appendArr) this.append(...appendArr)

    if (!result?.filter) return this

    Object.entries(result.filter).forEach(([key, value]) => {
      const values = Array.isArray(value) ? value : [value]
      this.createOrUpdateParams(`filter[${key}]`, values.map(String))
    })

    return this
  }

  toQueryParamResult(): QueryParamResult<Record<string, string | string[]>> {
    const result: QueryParamResult = {}

    // pagination
    const num = (k: string) => {
      const v = this._params[k]?.[0]
      const n = v !== undefined ? parseInt(v, 10) : NaN
      return Number.isNaN(n) ? undefined : n
    }

    result.page = num('page')
    result.per_page = num('per_page')
    result.limit = num('limit')

    // helper
    const singleOrArray = (vals?: string[]) => {
      if (!vals || vals.length === 0) return undefined
      return vals.length === 1 ? vals[0] : vals
    }

    if (this._params.sort) result.sort = singleOrArray(this._params.sort)
    if (this._params.include) result.include = singleOrArray(this._params.include)
    if (this._params.append) result.append = singleOrArray(this._params.append)

    // filters → keep `filter[...]` keys intact
    const filters: Record<string, string | string[]> = {}

    Object.entries(this._params).forEach(([key, vals]) => {
      const match = key.match(/^filter\[(.+)]$/)
      if (!match) return

      const field = match[1]
      const v = singleOrArray(vals)
      if (field !== undefined && v !== undefined) filters[field] = v
    })

    if (Object.keys(filters).length > 0) {
      result.filter = filters
    }

    return result
  }

  buildParams(key: string) {
    const param = this._params[key]
    if (param) {
      return `${key}=${param.join(',')}`
    }

    return ''
  }

  buildAsArray(options?: {
    includePath?: boolean
    includePagination?: boolean
    excludes?: string[]
  }) {
    const data: string[] = []

    if (options?.includePath) {
      data.push(this._path)
    }

    Object.keys(this._params).map((key) => {
      if (options?.excludes?.includes(key)) {
        return
      }

      if (
        !options?.includePagination &&
        (key === 'page' || key === 'per_page' || key === 'limit')
      ) {
        return
      }

      data.push(this.buildParams(key))
    })

    return data
  }

  build() {
    const params = Object.keys(this._params).map((key) => this.buildParams(key))

    const qMark = this.shouldHaveQmark() ? '?' : ''

    return `${this._path}${qMark}${this.connectQueryString(...params)}`
  }

  shouldHaveQmark() {
    return Object.keys(this._params).length > 0
  }

  connectQueryString(...strings: string[]) {
    return strings.filter((string) => string.length > 0).join('&')
  }
}

export default QueryBuilder
