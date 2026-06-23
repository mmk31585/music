/**
 * Persian (Farsi) translations for backend messages.
 *
 * The backend returns English messages in API responses.
 * This map translates them to Persian before showing to the user.
 *
 * Usage:
 *   import { translateMessage } from '@/utils/message-translations'
 *   toast.error(translateMessage('Invalid refresh token'))
 */

const messageMap: Record<string, string> = {
  // ── Auth ──────────────────────────────────────────────────────────
  'invalid refresh token': 'دوباره وارد شوید',
  'token expired': 'نشست شما منقضی شده است، دوباره وارد شوید',
  'invalid credentials': 'نام کاربری یا رمز عبور اشتباه است',
  'unauthorized': 'شما دسترسی ندارید',
  'forbidden': 'شما مجوز این کار را ندارید',
  'user not found': 'کاربر یافت نشد',
  'email already exists': 'این ایمیل قبلاً ثبت شده است',
  'username already taken': 'این نام کاربری قبلاً انتخاب شده است',
  'registration disabled': 'ثبت‌نام فعلاً غیرفعال است',

  // ── Tracks ────────────────────────────────────────────────────────
  'track not found': 'آهنگ یافت نشد',
  'track audio not found': 'فایل صوتی این آهنگ یافت نشد',
  'track is private': 'این آهنگ خصوصی است',
  'track created successfully': 'آهنگ با موفقیت ساخته شد',
  'track updated successfully': 'آهنگ با موفقیت به‌روزرسانی شد',
  'track deleted successfully': 'آهنگ با موفقیت حذف شد',

  // ── Albums ────────────────────────────────────────────────────────
  'album not found': 'آلبوم یافت نشد',
  'album created successfully': 'آلبوم با موفقیت ساخته شد',
  'album updated successfully': 'آلبوم با موفقیت به‌روزرسانی شد',
  'album deleted successfully': 'آلبوم با موفقیت حذف شد',

  // ── Artists ───────────────────────────────────────────────────────
  'artist not found': 'هنرمند یافت نشد',
  'artist created successfully': 'هنرمند با موفقیت ساخته شد',
  'artist updated successfully': 'هنرمند با موفقیت به‌روزرسانی شد',
  'artist deleted successfully': 'هنرمند با موفقیت حذف شد',

  // ── Playlists ─────────────────────────────────────────────────────
  'playlist not found': 'لیست پخش یافت نشد',
  'playlist created successfully': 'لیست پخش با موفقیت ساخته شد',
  'playlist updated successfully': 'لیست پخش با موفقیت به‌روزرسانی شد',
  'playlist deleted successfully': 'لیست پخش با موفقیت حذف شد',

  // ── Lyrics ────────────────────────────────────────────────────────
  'lyrics not found': 'متن آهنگ یافت نشد',
  'lyrics created successfully': 'متن آهنگ با موفقیت ذخیره شد',
  'lyrics updated successfully': 'متن آهنگ با موفقیت به‌روزرسانی شد',
  'lyrics deleted successfully': 'متن آهنگ با موفقیت حذف شد',
  'lyrics already exist': 'متن آهنگ قبلاً برای این زبان ثبت شده است',

  // ── Media / Upload ────────────────────────────────────────────────
  'file not found': 'فایل یافت نشد',
  'invalid file type': 'نوع فایل مجاز نیست',
  'file too large': 'حجم فایل بیش از حد مجاز است',
  'upload failed': 'آپلود با خطا مواجه شد',

  // ── Ingestion ─────────────────────────────────────────────────────
  'ingestion not found': 'درخواست ورود یافت نشد',
  'ingestion approved successfully': 'درخواست ورود با موفقیت تأیید شد',
  'ingestion rejected successfully': 'درخواست ورود رد شد',

  // ── Validation ────────────────────────────────────────────────────
  'invalid input': 'ورودی نامعتبر است',
  'invalid track id': 'شناسه آهنگ نامعتبر است',
  'invalid request body': 'درخواست نامعتبر است',
  'invalid media url': 'نشانی رسانه نامعتبر است',

  // ── Social ────────────────────────────────────────────────────────
  'party not found': 'گوش دادن گروهی یافت نشد',
  'room not found': 'اتاق یافت نشد',
  'club not found': 'باشگاه یافت نشد',
  'already a member': 'شما قبلاً عضو هستید',
  'not a member': 'شما عضو نیستید',

  // ── Generic ───────────────────────────────────────────────────────
  'internal server error': 'خطای داخلی سرور',
  'not found': 'یافت نشد',
  'request cancelled': 'درخواست لغو شد',
  'request timeout': 'زمان درخواست به پایان رسید',
  'too many requests': 'تعداد درخواست‌ها بیش از حد مجاز است',
}

/**
 * Translate an English message to Persian (Farsi).
 * Falls back to the original message if no translation exists.
 */
export function translateMessage(msg: string | null | undefined): string | null | undefined {
  if (!msg) return msg
  return messageMap[msg.toLowerCase()] ?? msg
}
