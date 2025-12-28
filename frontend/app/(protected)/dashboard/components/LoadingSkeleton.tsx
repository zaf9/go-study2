import React from 'react'
import { Skeleton } from 'antd'

export default function LoadingSkeleton() {
  return (
    <div style={{ display: 'grid', gap: 12 }}>
      <Skeleton active paragraph={{ rows: 1 }} title={{ width: '60%' }} />
      <div style={{ display: 'flex', gap: 12 }}>
        <Skeleton active paragraph={{ rows: 2 }} title={{ width: '100%' }} />
        <Skeleton active paragraph={{ rows: 2 }} title={{ width: '100%' }} />
        <Skeleton active paragraph={{ rows: 2 }} title={{ width: '100%' }} />
      </div>
      <Skeleton active paragraph={{ rows: 4 }} title={{ width: '100%' }} />
    </div>
  )
}
