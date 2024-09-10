import { Injectable } from '@nestjs/common';
import { ClickHouseClient, createClient } from '@clickhouse/client';
import { InjectDiscordClient } from '@discord-nestjs/core';
import { Client } from 'discord.js';
interface User {
  guildId: string;
  username: string;
  userId: string;
}
@Injectable()
export class MessageLog {
  private clickhouse: ClickHouseClient;

  constructor(
    @InjectDiscordClient()
    private readonly client: Client,
  ) {
    this.clickhouse = createClient({
      url: 'http://localhost:10100',
      database: 'discord',
    });
  }

  async insertMessageLog(data: any) {
    await this.clickhouse.insert({
      table: 'message',
      format: 'JSONEachRow',
      values: [data],
    });
  }

  async fetchFilteredAvatars(guildId: string, memberList: Array<unknown>) {
    try {
      const guild = await this.client.guilds.fetch(guildId);
      const members = await guild.members.fetch();

      const filteredMembers = members.filter((member) =>
        memberList.includes(member.user.id),
      );

      const avatarURLs = {};
      filteredMembers.forEach((member) => {
        avatarURLs[member.user.id] = member.user.avatarURL();
      });

      return avatarURLs;
    } catch (error) {
      console.error('Error fetching members or filtering:', error);
      return null;
    }
  }

  async messageCountForChannel(guildId: string) {
    const query = `
        SELECT channelId,
               channelName,
               COUNT(*) AS messageCount
        FROM discord.message
        WHERE guildId = '${guildId}'
          AND eventTime >= now() - toIntervalWeek(7)
        GROUP BY channelId, channelName
        ORDER BY messageCount DESC;
    `;
    const rawData = await this.clickhouse.query({
      query: query,
      format: 'JSON',
    });
    const data = await rawData.json();
    return data.data;
  }

  async sendMessageUsers(guildId: string): Promise<string[]> {
    const query = `
        SELECT DISTINCT ON (userId) username,
                                    userId
        FROM discord.message
        WHERE guildId = '${guildId}';
    `;
    const rawData = await this.clickhouse.query({
      query: query,
      format: 'JSON',
    });
    const data = await rawData.json();
    return data.data.map((user: User) => user.userId);
  }

  async messageCountForUsers(guildId: string) {
    const query = `
        SELECT userId,
               username,
               COUNT(*) AS messageCount
        FROM discord.message
        WHERE guildId = '${guildId}'
          AND eventTime >= now() - toIntervalDay(7) -- Son 7 gün
        GROUP BY userId, username
        ORDER BY messageCount DESC;
    `;
    const rawData = await this.clickhouse.query({
      query: query,
      format: 'JSON',
    });
    const result = await rawData.json();
    return result.data;
  }

  async messageChart(guildId: string) {
    const query = `
        SELECT toDate(eventTime) as eventDate,
               COUNT(*)          AS messageCount
        FROM discord.message
        WHERE guildId = '${guildId}'
          AND eventTime >= now() - toIntervalWeek(7)
        GROUP BY toDate(eventTime)
        ORDER BY eventDate;
    `;
    const rawData = await this.clickhouse.query({
      query: query,
      format: 'JSON',
    });
    const result = await rawData.json();
    return result.data;
  }

  async guildDashboard(guildId: string) {
    return {
      channel: await this.messageCountForChannel(guildId),
      users: await this.messageCountForUsers(guildId),
      chart: await this.messageChart(guildId),
      usersAvatarURL: await this.fetchFilteredAvatars(
        guildId,
        await this.sendMessageUsers(guildId),
      ),
    };
  }
}
