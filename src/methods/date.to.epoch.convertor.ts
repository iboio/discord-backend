export async function DateToEpochConvertor(eventTime: any) {
  const eventDateTime = new Date(eventTime.replace(' ', 'T'));
  return Math.floor(eventDateTime.getTime() / 1000);
}
