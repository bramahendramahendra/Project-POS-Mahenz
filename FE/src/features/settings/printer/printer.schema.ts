import { z } from 'zod'

export const printerSettingsSchema = z.object({
  paper_size: z.enum(['58mm', '80mm']),
  receipt_header: z.string(),
  receipt_footer: z.string(),
  show_logo: z.boolean(),
  auto_print: z.boolean(),
})

export type PrinterSettingsFormValues = z.infer<typeof printerSettingsSchema>
