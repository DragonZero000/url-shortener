const API_BASE_URL = ''

export interface ShortenResponse {
  url: string
}

export async function shortenUrl(url: string): Promise<ShortenResponse> {
  const response = await fetch(`${API_BASE_URL}/shorten`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ url }),
  })

  if (!response.ok) {
    const text = await response.text()
    throw new Error(text || `Server error: ${response.status}`)
  }

  return await response.json() as ShortenResponse
}