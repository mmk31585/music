import QueryBuilder from './query-builder'

const query = (path: string = '/') => {
  return new QueryBuilder(path)
}

export { query, QueryBuilder }
