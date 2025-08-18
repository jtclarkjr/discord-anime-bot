import { ChatInputCommandInteraction } from 'discord.js'
import { animeCommandDefinition } from '@/config/commands'

function formatOption(option: unknown, parent = '/anime'): string[] {
  const opt = option as { name: string; type: number; description: string; options?: unknown[] }
  let lines: string[] = []
  const fullName = `${parent} ${opt.name}`
  if (opt.type === 2 && opt.options) {
    // SUB_COMMAND_GROUP
    lines.push(`**${fullName}**: ${opt.description}`)
    for (const sub of opt.options) {
      lines = lines.concat(formatOption(sub, fullName))
    }
  } else if (opt.type === 1) {
    // SUB_COMMAND
    let usage = fullName
    if (opt.options && opt.options.length > 0) {
      const args = opt.options
        .map((o: unknown) => {
          const subOpt = o as { name: string; required?: boolean }
          const req = subOpt.required ? '<' : '['
          const end = subOpt.required ? '>' : ']'
          return `${req}${subOpt.name}${end}`
        })
        .join(' ')
      usage += ' ' + args
    }
    lines.push(`**${usage}**: ${opt.description}`)
  }
  return lines
}

export async function handleHelpCommand(interaction: ChatInputCommandInteraction) {
  const options = animeCommandDefinition.options
  let helpLines: string[] = []
  for (const opt of options) {
    helpLines = helpLines.concat(formatOption(opt))
  }
  const helpText = helpLines.join('\n')
  await interaction.reply({
    content: `Here are the available /anime commands:\n\n${helpText}`,
    flags: 1 << 6 // 64, EPHEMERAL flag
  })
}
