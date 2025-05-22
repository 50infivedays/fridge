'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

interface FridgeItem {
  item: string;
  quantity: number;
  unit: string;
  expireDate: string;
}

export default function SmartExpiry() {
  const [description, setDescription] = useState('');
  const [items, setItems] = useState<FridgeItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!description.trim()) {
      setError('请输入物品描述');
      return;
    }

    setLoading(true);
    setError('');
    
    try {
      const currentTime = new Date().toISOString().replace('T', ' ').substring(0, 19);
      
      const response = await fetch('http://localhost:8602/record', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          description,
          currentTime,
        }),
      });

      if (!response.ok) {
        throw new Error(`服务器返回错误: ${response.status}`);
      }

      const data = await response.json();
      setItems(data.items || []);
    } catch (err) {
      setError(`处理请求时出错: ${err instanceof Error ? err.message : String(err)}`);
    } finally {
      setLoading(false);
    }
  };

  // Format date to be more readable
  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      });
    } catch {
      return dateString; // Return original if parsing fails
    }
  };

  return (
    <div className="container mx-auto py-8 px-4">
      <h1 className="text-3xl font-bold mb-8 text-center">SmartExpiry - 智能物品管理</h1>
      
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>添加物品</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <label htmlFor="description" className="text-sm font-medium">
                物品描述（使用自然语言）
              </label>
              <Input
                id="description"
                placeholder="例如：2盒牛奶，保质期5天；3个苹果，下周一过期"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="w-full"
              />
              <p className="text-xs text-gray-500">
                提示：您可以一次性输入多个物品，并指定数量和过期时间
              </p>
            </div>
            <Button type="submit" disabled={loading} className="w-full">
              {loading ? '处理中...' : '添加到清单'}
            </Button>
            {error && <p className="text-red-500 text-sm">{error}</p>}
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>物品清单</CardTitle>
        </CardHeader>
        <CardContent>
          {items.length > 0 ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>物品名称</TableHead>
                  <TableHead>数量</TableHead>
                  <TableHead>单位</TableHead>
                  <TableHead>过期时间</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {items.map((item, index) => (
                  <TableRow key={index}>
                    <TableCell className="font-medium">{item.item}</TableCell>
                    <TableCell>{item.quantity}</TableCell>
                    <TableCell>{item.unit}</TableCell>
                    <TableCell>{formatDate(item.expireDate)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : (
            <p className="text-center py-4 text-gray-500">暂无物品，请添加物品</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
