import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { KPButton } from './KPButton'

describe('KPButton', () => {
  it('renders label', () => {
    render(<KPButton>Restart deployment</KPButton>)
    expect(screen.getByRole('button', { name: 'Restart deployment' })).toBeInTheDocument()
  })
})
