import { useNavigate } from "react-router-dom"
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { favoritesApi } from "@/api"
import DishImage from "@/components/DishImage"
import PageHeader from "@/components/PageHeader"
import type { Dish } from "@/types"
import toast from "react-hot-toast"
import { Heart } from "lucide-react"

export default function Favorites() {
  const navigate = useNavigate()
  const qc = useQueryClient()

  const { data: rawFavDishes, isLoading } = useQuery({
    queryKey: ["favorites"],
    queryFn: () => favoritesApi.list(),
  })
  const favDishes = rawFavDishes ?? []

  const favMut = useMutation({
    mutationFn: (dishId: number) => favoritesApi.remove(dishId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["favorites"] })
      qc.invalidateQueries({ queryKey: ["achievements"] })
      toast("已取消收藏")
    },
  })

  return (
    <div className="animate-fadeUp">
      <PageHeader title="收藏" subtitle="常吃和想吃的菜都在这里" icon={Heart} />

      <div className="px-5 py-4 max-w-[640px] mx-auto">
        {isLoading ? (
          <div className="grid grid-cols-2 gap-3">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="rounded-2xl overflow-hidden">
                <div className="skeleton h-[120px]" />
                <div className="p-2.5 space-y-2"><div className="skeleton h-3 w-3/4 rounded" /><div className="skeleton h-2.5 w-1/2 rounded" /></div>
              </div>
            ))}
          </div>
        ) : favDishes.length === 0 ? (
          <div className="text-center py-12">
            <span className="text-[56px] block mb-4 animate-float">❤</span>
            <div className="text-base font-semibold mb-1.5">还没有收藏</div>
            <div className="text-[13px] text-text2 mb-5">去看看有什么好吃的~</div>
            <button onClick={() => navigate("/dishes")} className="inline-flex items-center justify-center gap-1.5 py-2.5 px-5 rounded-full text-sm font-semibold bg-primary text-white">浏览菜品</button>
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-3">
            {favDishes.map((d: Dish) => (
              <div key={d.id} onClick={() => navigate(`/dishes/${d.id}`)} className="bg-card rounded-2xl overflow-hidden shadow-[0_1px_3px_rgba(0,0,0,.04),0_4px_12px_rgba(0,0,0,.04)] border border-border transition-all active:scale-97 cursor-pointer">
                <div className="h-[120px] bg-gradient-to-br from-primary-light to-pink-light relative">
                  <DishImage dish={d} className="w-full h-full" />
                  <button onClick={(e) => { e.stopPropagation(); favMut.mutate(d.id) }} className="absolute top-2 right-2 w-6 h-6 rounded-full bg-white/90 flex items-center justify-center text-xs text-primary active:animate-heartbeat">❤</button>
                </div>
                <div className="px-3 py-2.5">
                  <div className="text-sm font-semibold truncate mb-1">{d.name}</div>
                  <div className="flex gap-1.5 items-center text-[11px] text-text2">{d.category} <span className="w-[3px] h-[3px] rounded-full bg-text3" /> {d.cook_time}分钟</div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
