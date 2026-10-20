import { useNavigate } from 'react-router-dom'
import { FastConsole } from '../fast/FastConsole'

export function FastModePage() {
  const navigate = useNavigate()
  return (
    <div style={{ height: '100%' }}>
      <FastConsole onVisual={() => navigate('/visual')} />
    </div>
  )
}
