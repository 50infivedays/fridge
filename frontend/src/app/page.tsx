import Link from "next/link";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function Home() {
  return (
    <div className="container mx-auto py-12 px-4">
      <div className="flex flex-col items-center justify-center text-center space-y-8 mb-12">
        <h1 className="text-4xl font-bold tracking-tight">Fridge Manager</h1>
        <p className="text-xl text-gray-500 max-w-2xl">
          智能管理您的冰箱物品，追踪过期时间，避免食品浪费
        </p>
      </div>

      <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle>SmartExpiry</CardTitle>
            <CardDescription>智能物品管理</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-gray-500">
              使用自然语言添加物品，系统自动识别物品名称、数量和过期时间，帮助您轻松管理冰箱中的物品。
            </p>
          </CardContent>
          <CardFooter>
            <Link href="/smart-expiry" className="w-full">
              <Button className="w-full">立即使用</Button>
            </Link>
          </CardFooter>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>即将过期</CardTitle>
            <CardDescription>过期提醒</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-gray-500">
              查看即将过期的物品，避免食品浪费，合理安排食材使用顺序。
            </p>
          </CardContent>
          <CardFooter>
            <Button className="w-full" variant="outline" disabled>
              即将推出
            </Button>
          </CardFooter>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>购物清单</CardTitle>
            <CardDescription>智能推荐</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-gray-500">
              根据冰箱中的物品和使用习惯，智能生成购物清单，提醒您需要购买的物品。
            </p>
          </CardContent>
          <CardFooter>
            <Button className="w-full" variant="outline" disabled>
              即将推出
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
