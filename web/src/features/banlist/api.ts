import { banListApi } from '@/services/api';
import type { RequestsBanRequest } from '@appointment-master/api-client';

export function getBanList() {
  return banListApi.getBanList();
}

export function addToBanList(email: string) {
  return banListApi.addToBanList({ banRequest: { email } as RequestsBanRequest });
}

export function removeFromBanList(email: string) {
  return banListApi.removeFromBanList({ banRequest: { email } as RequestsBanRequest });
}

