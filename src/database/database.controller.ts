import { Controller, Get, Param } from '@nestjs/common';
import { DummyDataService } from './dummy.data.service';
import { MessageLog } from './log-service/discord/message.log';
import { DiscordGuild } from './log-service/discord/discord.guild';

@Controller('database')
export class DatabaseController {
  constructor(
    private readonly DCMessageLogService: MessageLog,
    private readonly discordGuild: DiscordGuild,
  ) {}

  @Get('/dashboard/:id')
  async guildDashboard(@Param('id') id: string) {
    return this.DCMessageLogService.guildDashboard(id);
  }
  @Get('/guilds')
  async guilds() {
    return this.discordGuild.getGuildList();
  }
}
