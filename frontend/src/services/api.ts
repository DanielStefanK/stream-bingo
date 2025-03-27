export interface Response<Data> {
  error: ErrorResponse
  success: boolean
  data: Data
}

export interface ErrorResponse {
  code: string
  msg: string
  data: Record<string, unknown>
}

// custom api error with code, msg and data

export class ApiError extends Error {
  code: string
  data: Record<string, unknown>

  constructor(code: string, msg: string, data: Record<string, unknown>) {
    super(msg)
    this.code = code
    this.data = data
  }
}

export async function fetchJson<T>(
  url: string,
  options?: RequestInit
): Promise<T> {
  const token = localStorage.getItem('token')
  const headers = new Headers(options?.headers)

  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(url, { ...options, headers })
  let resp: Response<T> | undefined = undefined
  try {
    resp = (await response.json()) as Response<T>
  } catch (error) {
    console.log(error)
    throw new Error('An error occurred while fetching the data.')
  }
  if (!resp.success) {
    if (resp.error && resp.error.code === 'AuthenticatingUser') {
      localStorage.removeItem('token')
    }
    throw new ApiError(resp.error.code, resp.error.msg, resp.error.data)
  }

  return resp.data
}

export interface PaginationResponse<T> {
  total: number
  data: T
}
