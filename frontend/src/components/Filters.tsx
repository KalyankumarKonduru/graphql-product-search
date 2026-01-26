import { useState } from 'react';
import { ChevronDown, SlidersHorizontal, X } from 'lucide-react';
import { useCategories } from '../hooks/useProducts';
import type {
  ProductFilterInput,
  ProductSortInput,
  ProductSortField,
  SortOrder,
} from '../graphql/types';

interface FiltersProps {
  filter: ProductFilterInput;
  sort: ProductSortInput;
  onFilterChange: (filter: ProductFilterInput) => void;
  onSortChange: (sort: ProductSortInput) => void;
}

const sortOptions: { label: string; field: ProductSortField; order: SortOrder }[] = [
  { label: 'Newest', field: 'CREATED_AT' as ProductSortField, order: 'DESC' as SortOrder },
  { label: 'Price: Low to High', field: 'PRICE' as ProductSortField, order: 'ASC' as SortOrder },
  { label: 'Price: High to Low', field: 'PRICE' as ProductSortField, order: 'DESC' as SortOrder },
  { label: 'Top Rated', field: 'RATING' as ProductSortField, order: 'DESC' as SortOrder },
  { label: 'Name: A-Z', field: 'NAME' as ProductSortField, order: 'ASC' as SortOrder },
];

const priceRanges = [
  { label: 'Under $50', min: 0, max: 50 },
  { label: '$50 - $100', min: 50, max: 100 },
  { label: '$100 - $500', min: 100, max: 500 },
  { label: '$500 - $1000', min: 500, max: 1000 },
  { label: 'Over $1000', min: 1000, max: undefined },
];

const ratingOptions = [4, 3, 2, 1];

export function Filters({
  filter,
  sort,
  onFilterChange,
  onSortChange,
}: FiltersProps) {
  const { categories, loading: categoriesLoading } = useCategories();
  const [showMobileFilters, setShowMobileFilters] = useState(false);

  const activeFiltersCount = [
    filter.categoryId,
    filter.minPrice !== undefined || filter.maxPrice !== undefined,
    filter.minRating,
    filter.inStock,
  ].filter(Boolean).length;

  const clearFilters = () => {
    onFilterChange({});
  };

  const handleCategoryChange = (categoryId: string | undefined) => {
    onFilterChange({ ...filter, categoryId });
  };

  const handlePriceChange = (min?: number, max?: number) => {
    onFilterChange({ ...filter, minPrice: min, maxPrice: max });
  };

  const handleRatingChange = (minRating?: number) => {
    onFilterChange({ ...filter, minRating });
  };

  const handleStockChange = (inStock?: boolean) => {
    onFilterChange({ ...filter, inStock });
  };

  const handleSortChange = (field: ProductSortField, order: SortOrder) => {
    onSortChange({ field, order });
  };

  const currentSortLabel =
    sortOptions.find(
      (o) => o.field === sort.field && o.order === sort.order
    )?.label || 'Sort';

  return (
    <>
      {/* Desktop Filters */}
      <div className="hidden lg:flex items-center gap-4 flex-wrap">
        {/* Sort Dropdown */}
        <div className="relative">
          <select
            value={`${sort.field}-${sort.order}`}
            onChange={(e) => {
              const [field, order] = e.target.value.split('-') as [ProductSortField, SortOrder];
              handleSortChange(field, order);
            }}
            className="appearance-none bg-white border border-gray-300 rounded-lg px-4 py-2 pr-10 text-sm font-medium text-gray-700 focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
            {sortOptions.map((option) => (
              <option key={`${option.field}-${option.order}`} value={`${option.field}-${option.order}`}>
                {option.label}
              </option>
            ))}
          </select>
          <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
        </div>

        {/* Category Filter */}
        <div className="relative">
          <select
            value={filter.categoryId || ''}
            onChange={(e) => handleCategoryChange(e.target.value || undefined)}
            disabled={categoriesLoading}
            className="appearance-none bg-white border border-gray-300 rounded-lg px-4 py-2 pr-10 text-sm font-medium text-gray-700 focus:outline-none focus:ring-2 focus:ring-primary-500 disabled:opacity-50"
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

        {/* Price Filter */}
        <div className="relative">
          <select
            value={
              filter.minPrice !== undefined
                ? `${filter.minPrice}-${filter.maxPrice || ''}`
                : ''
            }
            onChange={(e) => {
              if (!e.target.value) {
                handlePriceChange(undefined, undefined);
              } else {
                const [min, max] = e.target.value.split('-').map(Number);
                handlePriceChange(min, max || undefined);
              }
            }}
            className="appearance-none bg-white border border-gray-300 rounded-lg px-4 py-2 pr-10 text-sm font-medium text-gray-700 focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
            <option value="">Any Price</option>
            {priceRanges.map((range) => (
              <option
                key={range.label}
                value={`${range.min}-${range.max || ''}`}
              >
                {range.label}
              </option>
            ))}
          </select>
          <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
        </div>

        {/* Rating Filter */}
        <div className="relative">
          <select
            value={filter.minRating || ''}
            onChange={(e) =>
              handleRatingChange(e.target.value ? Number(e.target.value) : undefined)
            }
            className="appearance-none bg-white border border-gray-300 rounded-lg px-4 py-2 pr-10 text-sm font-medium text-gray-700 focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
            <option value="">Any Rating</option>
            {ratingOptions.map((rating) => (
              <option key={rating} value={rating}>
                {rating}+ Stars
              </option>
            ))}
          </select>
          <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
        </div>

        {/* In Stock Toggle */}
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={filter.inStock || false}
            onChange={(e) => handleStockChange(e.target.checked || undefined)}
            className="w-4 h-4 text-primary-600 rounded focus:ring-primary-500"
          />
          <span className="text-sm font-medium text-gray-700">In Stock Only</span>
        </label>

        {/* Clear Filters */}
        {activeFiltersCount > 0 && (
          <button
            onClick={clearFilters}
            className="flex items-center gap-1 text-sm font-medium text-red-600 hover:text-red-700"
          >
            <X className="w-4 h-4" />
            Clear ({activeFiltersCount})
          </button>
        )}
      </div>

      {/* Mobile Filter Button */}
      <div className="lg:hidden">
        <button
          onClick={() => setShowMobileFilters(true)}
          className="flex items-center gap-2 px-4 py-2 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700"
        >
          <SlidersHorizontal className="w-4 h-4" />
          Filters
          {activeFiltersCount > 0 && (
            <span className="bg-primary-600 text-white text-xs px-2 py-0.5 rounded-full">
              {activeFiltersCount}
            </span>
          )}
        </button>
      </div>

      {/* Mobile Filters Modal */}
      {showMobileFilters && (
        <div className="fixed inset-0 z-50 lg:hidden">
          <div
            className="absolute inset-0 bg-black/50"
            onClick={() => setShowMobileFilters(false)}
          />
          <div className="absolute right-0 top-0 h-full w-full max-w-sm bg-white shadow-xl overflow-y-auto">
            <div className="p-4 border-b border-gray-200 flex items-center justify-between">
              <h2 className="text-lg font-semibold">Filters</h2>
              <button
                onClick={() => setShowMobileFilters(false)}
                className="p-2 hover:bg-gray-100 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <div className="p-4 space-y-6">
              {/* Mobile filter options would go here - similar to desktop */}
              <p className="text-gray-500 text-sm">
                Filter options for mobile view
              </p>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
