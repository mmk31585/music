import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import NotificationItem from '@/components/music/social/NotificationItem.vue'

vi.mock('primevue/usetoast', () => ({
  useToast: vi.fn(() => ({
    add: vi.fn(),
  })),
}))

function makeNotification(overrides: Record<string, unknown> = {}) {
  return {
    id: 'n1',
    type: 'user_followed',
    title: 'New Follower',
    body: 'Someone followed you',
    isRead: false,
    createdAt: '2024-06-15T10:00:00Z',
    data: null,
    ...overrides,
  }
}

describe('NotificationItem', () => {
  it('renders notification title', () => {
    const wrapper = mount(NotificationItem, {
      props: { notification: makeNotification({ title: 'Hello' }) },
    })
    expect(wrapper.text()).toContain('Hello')
  })

  it('renders notification message', () => {
    const wrapper = mount(NotificationItem, {
      props: { notification: makeNotification({ body: 'Welcome to the platform' }) },
    })
    expect(wrapper.text()).toContain('Welcome')
  })

  it('shows mark-read button for unread notifications', () => {
    const wrapper = mount(NotificationItem, {
      props: { notification: makeNotification({ isRead: false }) },
    })
    expect(wrapper.find('button').exists()).toBe(true)
  })
})
