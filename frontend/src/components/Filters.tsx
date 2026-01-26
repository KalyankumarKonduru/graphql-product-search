import { ChevronDown } from 'lucide-react';
import { useCategories } from '../hooks/useProducts';
import type { ProductFilterInput, ProductSortInput } from '../graphql/types';

interface FiltersProps {
  filter: ProductFilterInput;
  sort: ProductSortInput;
  onFilterChange: (filter: ProductFilterInput) => void;
  onSortChange: (sort: ProductSortInput) => void;
}

export function Filters({ filter, sort, onFilterChange, onSortChange }: FiltersProps) {
  const { categories } = useCategories();

  return (
    <div className="flex flex-wrap items-center gap-4">
      <div className="relative">
        <select
          value={`${sort.field}-${sort.order}`}
          onChange={(e) => {
            const [field, order] = e.target.value.split('-') as [ProductSortInput['field'], ProductSortInput['order']];
            onSortChange({ field, order });
          }}
          className="appearance-none bg-white border border-gray-300 rounded-lg px-4 py-2 pr-10 text-sm"
        >
          <option value="CREATED_AT-DESC">Newest</option>
          <option value="PRICE-ASC">Price: Low to High</option>
          <option value="PRICE-DESC">Price: High to Low</option>
          <option value="RATING-DESC">Top Rated</option>
          <option value="NAME-ASC">Name: A-Z</option>
        </select>
        <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
      </div>

      <div className="relative">
        <select
          value={filter.categoryId || ''}
          onChange={(e) => onFilterChange({ ...filter, categoryId: e.target.value || undefined })}
          className="appearance-none bg-white border border-gray-300 rounded-lg px-4 py-2 pr-10 text-sm"
        >
          <option value="">All Categories</option>
          {categories.map((cat: { id: string; name: string; productCount: number }) => (
            <option key={cat.id} value={cat.id}>
              {cat.name} ({cat.productCount})
            </option>
          ))}
        </select>
        <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
      </div>

      <label className="flex items-center gap-2 cursor-pointer">
        <input
          type="checkbox"
          checked={filter.inStock || false}
          onChange={(e) => onFilterChange({ ...filter, inStock: e.target.checked || undefined })}
          className="w-4 h-4 text-primary-600 rounded"
        />
        <span className="text-sm text-gray-700">In Stock Only</span>
      </label>
    </div>
  );
}
