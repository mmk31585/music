export const dateUtil = {
  toPersianDate: (timestamp: string | number | null | undefined) => {
    if (!timestamp) return '-'

    const date = new Date(timestamp)

    if (isNaN(date.valueOf())) return 'تاریخ نامعتبر'

    return `${date.toLocaleString('fa-IR', { day: 'numeric' })} ${date.toLocaleString('fa-IR', { month: 'long', calendar: 'persian' })} ${date.toLocaleString('fa-IR', { year: 'numeric', calendar: 'persian' })}`
  },

  toPersianTime: (timestamp: string | number | null | undefined) => {
    if (!timestamp) return '-'

    const date = new Date(timestamp)

    if (isNaN(date.valueOf())) return 'تاریخ نامعتبر'

    return date.toLocaleTimeString('fa-IR', {
      hour: 'numeric',
      minute: 'numeric',
      hour12: false,
      calendar: 'persian',
    })
  },

  toPersianDateTime: (timestamp: string | number | null | undefined) => {
    if (!timestamp) return '-'

    const date = new Date(timestamp)

    if (isNaN(date.valueOf())) return 'تاریخ نامعتبر'

    const dayMonth = date.toLocaleDateString('fa-IR', {
      day: 'numeric',
      month: 'long',
      calendar: 'persian',
    })

    const year = date.toLocaleDateString('fa-IR', {
      year: 'numeric',
      calendar: 'persian',
    })

    const time = date.toLocaleTimeString('fa-IR', {
      hour: 'numeric',
      minute: 'numeric',
      hour12: false,
      calendar: 'persian',
    })

    return `${dayMonth} ${year} ساعت ${time}`
  },
}
