export function formatBytes(bytes) {
  if (bytes < 1024 * 1024) {
    return (bytes / 1024).toFixed(2) + ' KB'
  } else {
    return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
  }
}

export function getThumbFilepath (filepath) {
  if (!filepath) {
    return filepath
  }

  const [pathPart, queryPart] = filepath.split('?')
  const urlParts = pathPart.split('/')
  const filename = urlParts.pop()
  urlParts.push(`thumb_${filename}`)

  const thumbPath = urlParts.join('/')
  return queryPart ? `${thumbPath}?${queryPart}` : thumbPath
}
