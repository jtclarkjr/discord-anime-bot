import { watchlistFile } from '@config/constants'

async function readWatchlists(): Promise<Record<string, number[]>> {
  try {
    const file = Bun.file(watchlistFile)
    if (!(await file.exists())) return {}
    const raw = await file.text()
    return JSON.parse(raw)
  } catch {
    return {}
  }
}

async function writeWatchlists(data: Record<string, number[]>) {
  await Bun.write(watchlistFile, JSON.stringify(data, null, 2))
}

export async function getUserWatchlist(userId: string): Promise<number[]> {
  const lists = await readWatchlists()
  return lists[userId] || []
}

export async function addToWatchlist(
  userId: string,
  animeId: number
): Promise<{ success: boolean; message: string }> {
  const lists = await readWatchlists()
  if (!lists[userId]) lists[userId] = []
  if (lists[userId].includes(animeId)) {
    return { success: false, message: 'Anime already in your watchlist.' }
  }
  lists[userId].push(animeId)
  await writeWatchlists(lists)
  return { success: true, message: 'Anime added to your watchlist.' }
}

export async function removeFromWatchlist(
  userId: string,
  animeId: number
): Promise<{ success: boolean; message: string }> {
  const lists = await readWatchlists()
  if (!lists[userId] || !lists[userId].includes(animeId)) {
    return { success: false, message: 'Anime not found in your watchlist.' }
  }
  lists[userId] = lists[userId].filter((id) => id !== animeId)
  await writeWatchlists(lists)
  return { success: true, message: 'Anime removed from your watchlist.' }
}
