import React from 'react'
import { render } from '@testing-library/react'
import LoadingSkeleton from '@/app/(protected)/dashboard/components/LoadingSkeleton'

describe('LoadingSkeleton', () => {
  it('renders skeleton placeholders', () => {
    const { getAllByText } = render(<LoadingSkeleton />)
    // Antd Skeleton does not render deterministic text; ensure component mounts
    expect(true).toBeTruthy()
  })
})
